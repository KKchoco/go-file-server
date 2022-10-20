function dl() {
    let a = passwordInput.value;
    if (!a) return alert("Please enter the upload password before downloading the ShareX config.");
    var b = document.createElement("a"),
        c = new Blob(
            [
                JSON.stringify({
                    Version: "13.7.0",
                    Name: `go-file-server [${wl}]`,
                    DestinationType: "ImageUploader, FileUploader",
                    RequestMethod: "POST",
                    RequestURL: `${wl}/upload`,
                    Body: "MultipartFormData",
                    Arguments: { password: a },
                    FileFormName: "file",
                    URL: "$json:url$",
                    DeletionURL: "$json:deletionUrl$",
                    ErrorMessage: "$json:error$",
                }),
            ],
            { type: "text/plain" }
        );
    (b.href = URL.createObjectURL(c)), (b.download = "go-file-server.sxcu"), b.click();
}
dl();