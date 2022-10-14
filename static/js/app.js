function Prompt() {
    let toast = function (c) {
        const{
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    async function custom(c) {
        const {
            icon = "",
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        const { value: results } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,

            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },

            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },

            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            },
        })

        if (results) {
            if (results.dismiss !== Swal.DismissReason.cancel ) {
                if (results.value !== "" ) {
                    if (c.callback !== undefined) {
                        c.callback(results);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }
    }  

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

// function roomIdForRoomTmpl (room_id) {
//     attention.custom({
//         title: 'Choose your dates',
//         msg: html,

//         willOpen: () => {
//             const elem = document.getElementById("reservation-dates-modal");
//             const rp = new DateRangePicker(elem, {
//                 format: 'yyyy-mm-dd',
//                 showOnFocus: true,
//                 minDate: new Date(),
//             })
//         },

//         didOpen: () => {
//             document.getElementById("start").removeAttribute("disabled");
//             document.getElementById("end").removeAttribute("disabled");
//         },
        
//         callback: function(results) {
//             console.log("called");
            
//             let form = document.getElementById("check-availability-form")
//             let formData = new FormData(form);
//             formData.append("csrf_token", "{{.CSRFToken}}");
//             formData.append("room_id", room_id);

//             fetch('/search-availability-json', {
//                 method:"post",
//                 body:formData,
                
//             })
//                 .then(response => response.json())
//                 .then(data => {
//                     if (data.ok) {
//                         attention.custom({
//                             icon: 'success',
//                             showConfirmButton: false,
//                             msg: '<p> Room is available!<p>'
//                                 + '<p><a href="/book-room?id='
//                                 + data.room_id
//                                 +'&s='
//                                 + data.start_date
//                                 +'&e='
//                                 + data.end_date
//                                 +'" class="btn btn-primary">'
//                                 + 'Book now !</a></p>',
//                         })
//                     } else {
//                         attention.error({
//                             msg: "No availabilty!",
//                         });
//                     }
//                 })
//             }
//         });
//     }