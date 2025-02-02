// Zip Implementation:
// It is slow in intial load

import JSZip from "jszip";
import { loadImage } from "./main";


export async function ziploadPreviewImages(): Promise<void> {
    const container = document.getElementById('container');
    //clearing the container
    if (!container) return;
    container.innerHTML = '';
    try {
        const res = await fetch('/api/images/all');
        if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);

        // Get the zip files
        const zipBlob = await res.blob();

        // Using JSZip to extract images
        const zip = await JSZip.loadAsync(zipBlob);

        const imagePromises = Object.keys(zip.files).map(async filename => {
            const file = zip.files[filename];
            if (!file.dir) {
                const blob = await file.async("blob");
                const imgURL = URL.createObjectURL(blob);

                const img = document.createElement("img");
                img.src = imgURL;
                img.id = filename.split("_")[1].split(".")[0];
                img.width = 160;
                img.height = 90;
                img.onload = () => URL.revokeObjectURL;
                return img;
            }
        })

        const images = await Promise.all(imagePromises);

        images.forEach(img => {
            if (img) container?.appendChild(img);
        })

        // click on image
        const imgTags = document.querySelectorAll("#container img");
        imgTags.forEach(img => img.addEventListener("click", () => loadImage(img)));

    } catch (err) {
        console.error('Error loading images:', err);
    }
}

