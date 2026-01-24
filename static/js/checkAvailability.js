function checkAvailability(roomId, csrfToken) {
    let html = `
    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
        <div class="form-row">
            <div class="col">
                <div class="form-row" id="reservation-dates-modal">
                    <div class="col mt-2">
                        <input disabled required class="form-control" type="text" name="start" id="start" 
                        placeholder="Arrival" autocomplete="off">
                    </div>
                    <div class="col mt-2">
                        <input disabled required class="form-control" type="text" name="end" id="end" 
                        placeholder="Departure" autocomplete="off">
                    </div>

                </div>
            </div>
        </div>
    </form>
    `

    attention.custom({
        title: 'Choose your dates',
        msg: html,

        //willOpen and didOpen are called in case.layout.tmpl
        willOpen: () => {
            const elem = document.getElementById('reservation-dates-modal');
            const rp = new DateRangePicker(elem, {
                format: 'dd-mm-yyyy',
                showOnFocus: true,
                minDate: new Date(),
            })
        },
        didOpen: () => {
            document.getElementById('start').removeAttribute('disabled');
            document.getElementById('end').removeAttribute('disabled')
        },

        callback: function (result) {
            console.log("called");

            //this variable is storing form
            let form = document.getElementById("check-availability-form");
            //this variable is storing all elements from the form
            let formData = new FormData(form);
            //this is appending CSRFToken from any template that we created
            formData.append("csrf_token", csrfToken)
            // we add room id, which it 1
            formData.append("room_id", roomId)

            fetch('/search-availability-json', {
                method: "post",
                body: formData,
            }).then(response => response.json())
                .then(data => {
                    if (data.ok) {
                        attention.custom({
                            icon: 'success',
                            showConfirmButton: false,
                            msg: '<p>Room is available!</p>'
                                + '<p><a href="book-room?id='
                                + data.room_id
                                + '&s='
                                + data.start_date
                                + '&e='
                                + data.end_date
                                + '" class="btn btn-primary">'
                                + 'Book now!</a></p>'
                        })
                    } else {
                        attention.error({
                            msg: "No availability"
                        })
                    }
                })
        }
    });
}
