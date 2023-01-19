// if (localStorage.getItem("uid") !== null) {
//     localStorage.removeItem("uid")
// }

const sub_btn = document.getElementById("sub_btn")
const form = document.getElementById("form")
const maxFileSize = 1024 * 1024
const uploadFiles = document.getElementById("myFile")
const filters = document.getElementById("filters")
let filter_values = document.getElementsByName("filter");
let err_div = document.getElementById("error")

let oldDownload;

(async function () {

    await pingServer();

    let uid = localStorage.getItem("uid")
    if (!uid) {
        return
    }

    oldDownload = await fetch(`http://localhost:5000/check/${uid}`)
    if (oldDownload.status !== 200) {
        return
    }

    oldDownload = document.getElementById("old_download");

    if (oldDownload) {
        document.getElementById("old_files_container").style.display = "block";
        oldDownload.addEventListener("click", (e) => {
            e.preventDefault();
            dowloadFiles();
        })
    }
})()



uploadFiles.addEventListener("change", async (e) => {
    let filesCount = uploadFiles.files.length
    if (filesCount > 0) {
        filters.style.display = "block";
    } else {
        filters.style.display = "none";
    }
})

sub_btn.addEventListener("click", async (e) => {
    e.preventDefault()

    let uid = crypto.randomUUID();
    const formData = new FormData()
    const file = document.getElementById("myFile")

    let files = arrayOfFiles(file.files)
    const totalSize = files.reduce((acc, elm) => elm.size + acc, 0)

    for (let i = 0; i < files.length; i++) {
        formData.append('files', files[i])
    }

    let filterValue = getFilterValue()
    if (filterValue === "") {
        return
    }


    formData.append('filter', filterValue)
    formData.append('uid', uid)


    try {
        let res = await fetch("http://localhost:5000/upload", { method: "POST", body: formData })
        if (res.status === 200) {
            err_div.innerHTML = ""
            window.location = "download.html"
            localStorage.setItem("uid", uid)
        } else {
            let err = await res.text()
            matchError(res.status, err)
        }
    } catch (err) {
        if (err.message === "Failed to fetch") {
            window.location = "serviceDown.html"
        }
    }
});

function matchError(statusCode, err) {
    switch (statusCode) {
        case 400:
            matchClientErrors(err)
            break;

        case 500:
            matchServerErrors(err)
            break;

        default:
            console.log("error is not matching")
    }
}

function matchClientErrors(err) {
    if (err === "File is too large") {
        let err_div = document.getElementById("error")
        err_div.innerHTML = err
        err_div.style.color = "red"
    }
}


function matchServerErrors(err) {
    return undefined
}



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

async function pingServer() {
    await fetch("http://localhost:5000/test").catch(() => window.location = "serviceDown.html")
}


async function dowloadFiles() {
    // send request to check if file is created with the uid
    // let res = await fetch(`http://localhost:5000/check/${uid}`)

    // if (res.status !== 200) {
    //     return
    // }
    let uid = localStorage.getItem("uid");
    res = await fetch(`http://localhost:5000/download/${uid}`)

    if (res.status !== 200) {
        return
    }

    const bloby = await res.blob()
    const href = URL.createObjectURL(bloby)
    const a = Object.assign(document.createElement('a'), { href, style: "display:none", download: "Images" })

    document.body.append(a)
    a.click()

    URL.revokeObjectURL(href)
    a.remove()

}

function getFilterValue() {
    for (let i = 0; i < filter_values.length; i++) {
        if (filter_values[i].checked) {
            return filter_values[i].value;
        }
    }
    return ""
}