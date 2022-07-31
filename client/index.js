const sub_btn = document.getElementById("sub_btn")
const form = document.getElementById("form")
const maxFileSize = 1024 * 1024

sub_btn.addEventListener("click", async (e) => {
    e.preventDefault()
    const formData = new FormData()
    const file = document.getElementById("myFile")

    let files = arrayOfFiles(file.files)
    const totalSize = files.reduce((acc, elm) => elm.size + acc, 0)

    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i])
    }
    formData.append('filter', "gray")

    console.log(formData)
    try {
        let res = await fetch("http://localhost:5000/upload", { method: "POST", body: formData })
        console.log(res)
    } catch (err) {
        if ((err.message === "Failed to fetch" && totalSize >= maxFileSize) || (res.status === 400)) {
            console.log("File is too large")
        }
    }
});

function arrayOfFiles(objectFiles) {
    let files = new Array();
    let keys = Object.keys(objectFiles)
    for (let i = 0; i < objectFiles.length; i++) {
        if (parseInt(keys[i]) !== NaN) {
            files.push(objectFiles[i])
        }
    }
    return files
}
