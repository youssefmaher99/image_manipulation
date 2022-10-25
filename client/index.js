if (localStorage.getItem("uid") !== null) {
    localStorage.removeItem("uid")
}

const sub_btn = document.getElementById("sub_btn")
const form = document.getElementById("form")
const maxFileSize = 1024 * 1024
const uid = crypto.randomUUID();

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
    formData.append('uid', uid)

    
    try {
        let res = await fetch("http://localhost:5000/upload", { method: "POST", body: formData })
        localStorage.setItem("uid", uid)
        if (res.status === 200) {
            // window.location = "file:///home/youssef/Desktop/sandbox/image_manipulation/client/download.html"
        }
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
