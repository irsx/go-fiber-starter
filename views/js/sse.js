if (typeof EventSource !== "undefined") {
    var source = new EventSource("/api/import/products/stream");
    source.onmessage = function (event) {
        if (!event) {
            source.close();
        } else {
            const index = parseInt($("#counter-progress").text()) + 1;
            $("#counter-progress").text(index);

            var alertStatus = "success";
            if (event.data.includes("error")) {
                alertStatus = "danger";
            }

            const alertResult = `<div class="alert alert-${alertStatus} mb-2" role="alert">#${index} ${event.data}</div>`;
            document.getElementById("result").innerHTML += alertResult;
        }
    };
    window.closedConnection = () => {
        source.close();
    };
    window.onbeforeunload = () => {
        console.log("closing");
        return "closing unload";
    };
    window.onunload = () => {
        console.log("closed");
        source.close();
    };
} else {
    document.getElementById("result").innerHTML =
        "Sorry, your browser does not support server-sent events...";
}

const send = () => {
    var input = document.querySelector('input[type="file"]');
    var data = new FormData();
    // replace jwt token below
    var jwtToken =
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3N1ZXIiOiJlZDAyODAxNi0yMjQ1LTQxYWEtOGMxZS0wYjJkNjI5OTBhZGYiLCJleHAiOjMzOTg0NTgwNjR9.8D39TRB9QtEwmUcbRRFFDbXVt_bPNujajmCphFdeS3Q";

    data.append("file", input.files[0]);

    fetch("/api/import/products", {
        method: "POST", // *GET, POST, PUT, DELETE, etc.
        mode: "cors", // no-cors, *cors, same-origin
        cache: "no-cache", // *default, no-cache, reload, force-cache, only-if-cached
        credentials: "same-origin", // include, *same-origin, omit
        redirect: "follow", // manual, *follow, error
        referrerPolicy: "no-referrer", // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
        body: data, // body data type must match "Content-Type" header
        headers: {
            Authorization: `Bearer ${jwtToken}`,
        },
    })
        .then((response) => response.json())
        .then((data) => {
            $("#counter-area").show();
            $("#counter-total").text(data.data.total);
        });
};
const closer = () => {
    source.close();
};
