// import JSZip from "jszip";

const display = document.getElementById("display-container")!;

async function loadPreviewImages() {
    const container = document.getElementById('container');
    if (!container) {
        console.error("Container not found!");
        return;
    }

    container.innerHTML = '';
    const promiseBuffer: Promise<void>[] = [];
    const resolution = "low";

    for (let i = 1; i <= 15; i++) {
        const img = document.createElement("img");
        img.id = i.toString();
        img.width = 160;
        img.height = 90;
        const fetchPromises = fetch(`/api/images/${resolution}/${img.id}`)
            .then(res => {
                if (!res.ok) {
                    throw new Error(`Response status: ${res.status}`);
                }
                return res.blob();

            }).then(blob => {
                img.src = URL.createObjectURL(blob)
                container?.appendChild(img);

            }).catch(err => {
                console.error(`Error loading image ${img.id}:`, err);
            });
        promiseBuffer.push(fetchPromises);
    }

    try {
        await Promise.all(promiseBuffer);

        // Add event listeners to all images
        const imgTags = document.querySelectorAll("#container img");
        imgTags.forEach((img) => {
            img.addEventListener("click", () => loadImage(img));
        });
    } catch (err) {
        console.error("Error loading images:", err);
    }
}

// Zip Implementation:
// async function ziploadPreviewImages(): Promise<void> {
//     //clearing the container
//     if (!container) return;
//     container.innerHTML = '';
//     try {
//         const res: Response = await fetch('/api/images/all');
//         if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
//
//         // Get the zip files
//         const zipBlob = await res.blob();
//
//         // Using JSZip to extract images
//         const zip = await JSZip.loadAsync(zipBlob);
//
//         const imagePromises = Object.keys(zip.files).map(async filename => {
//             const file = zip.files[filename];
//             if (!file.dir) {
//                 const blob = await file.async("blob");
//                 const imgURL = URL.createObjectURL(blob);
//
//                 const img = document.createElement("img");
//                 img.src = imgURL;
//                 img.id = filename.split("_")[1].split(".")[0];
//                 img.width = 160;
//                 img.height = 90;
//                 img.onload = () => URL.revokeObjectURL;
//                 return img;
//             }
//         })
//
//         const images = await Promise.all(imagePromises);
//
//         images.forEach(img => {
//             if (img) container?.appendChild(img);
//         })
//
//         // click on image
//         const imgTags = document.querySelectorAll("#container img");
//         imgTags.forEach(img => img.addEventListener("click", () => loadImage(img)));
//
//     } catch (err) {
//         console.error('Error loading images:', err);
//     }
// }

async function loadImage(imgElement: Element) {
    const displayImg = document.createElement("img");
    const resolution = "high";
    displayImg.width = 1280;
    displayImg.height = 720;
    try {
        const res = await fetch(`/api/images/${resolution}/${imgElement.id}`);
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const blob = await res.blob()
        displayImg.src = URL.createObjectURL(blob)
    } catch (err) {
        console.error(err);
    }
    display.innerHTML = '';
    display.appendChild(displayImg);
}


const btn = document.getElementById("btn")!;
btn.addEventListener("click", loadPreviewImages);
