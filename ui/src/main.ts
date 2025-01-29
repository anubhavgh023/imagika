const display = document.getElementById("display-container")!;

const TOTAL_IMAGES = 15

// async function loadRandomImages() {
//     const container = document.getElementById("container")!;
//     container.innerHTML = '';
//     for (let i = 1; i <= 15; i++) {
//         const img = document.createElement("img");
//         img.id = i.toString();
//         const resolution = "low";
//         try {
//             const res = await fetch(`/api/images/${resolution}/${img.id}`);
//             if (!res.ok) {
//                 throw new Error(`Response status: ${res.status}`);
//             }
//
//             const blob = await res.blob()
//             img.src = URL.createObjectURL(blob)
//
//         } catch (err) {
//             console.error(err);
//         }
//         img.width = 160;
//         img.height = 90;
//         container.appendChild(img);
//     }
//     const imgTags = document.querySelectorAll("#container img");
//     imgTags.forEach(img => img.addEventListener("click", () => loadImage(img)));
// }

async function loadPreviewImages(): Promise<void> {
    const container: HTMLElement | null = document.getElementById('container');
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 seconds


    try {
        const res: Response = await fetch('/api/images/all', { signal: controller.signal });
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }

        const reader: ReadableStreamDefaultReader<Uint8Array> | undefined = res.body?.getReader();
        if (!reader) {
            throw new Error('Failed to get reader from response body');
        }

        let imagesReceived: number = 0;

        while (imagesReceived < TOTAL_IMAGES) { // Assuming 15 images
            // Read the image size (4 bytes)
            const sizeBuffer = await reader.read();
            console.log("sizeBuffer:", sizeBuffer);
            if (sizeBuffer.done || !sizeBuffer.value) break;


            // Extract the size from the first 4 bytes
            const size: number = new DataView(sizeBuffer.value.buffer).getUint32(0, false);
            console.log("size:", size);

            // Read the image data
            let imgData: Uint8Array = new Uint8Array(size);
            console.log("imgData:", imgData);
            let bytesRead = 0;

            // Issue Here
            while (bytesRead < size) {
                const chunk: ReadableStreamReadResult<Uint8Array> = await reader.read();
                console.log("chunck:", chunk);
                if (chunk.done || !chunk.value) break;

                // Append the chunk to the image data
                imgData.set(new Uint8Array(chunk.value), bytesRead);
                console.log("imgData:", imgData);
                bytesRead += chunk.value.length;
            }

            if (bytesRead !== size) {
                throw new Error(`Incomplete image data: expected ${size} bytes, but got ${bytesRead}`);
            }

            console.log(imgData);

            // Create an image element and display it
            const blob: Blob = new Blob([imgData], { type: 'image/png' });
            console.log("blob:", blob);
            const imgURL: string = URL.createObjectURL(blob);
            console.log("imgURL:", imgURL);
            const img: HTMLImageElement = document.createElement('img');
            img.src = imgURL;
            img.width = 160;
            img.height = 90;
            container?.appendChild(img);

            imagesReceived++;
        }
    } catch (err) {
        console.error('Error loading images:', err);
    }
}

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
