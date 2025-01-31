// import JSZip from "jszip";

const display = document.getElementById("display-container")!;
const container: HTMLElement | null = document.getElementById('container');

async function loadPreviewImages() {
    if (container) {
        container.innerHTML = '';
    }
    for (let i = 1; i <= 15; i++) {
        const img = document.createElement("img");
        img.id = i.toString();
        const resolution = "low";
        try {
            const res = await fetch(`/api/images/${resolution}/${img.id}`);
            if (!res.ok) {
                throw new Error(`Response status: ${res.status}`);
            }

            const blob = await res.blob()
            img.src = URL.createObjectURL(blob)

        } catch (err) {
            console.error(err);
        }
        img.width = 160;
        img.height = 90;
        container?.appendChild(img);
    }
    const imgTags = document.querySelectorAll("#container img");
    imgTags.forEach(img => img.addEventListener("click", () => loadImage(img)));
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
