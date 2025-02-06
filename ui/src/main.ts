const TOTAL_IMAGES = 28

async function loadPreviewImages() {
    const container = document.getElementById('container');
    if (!container) {
        console.error("Container not found!");
        return;
    }

    container.innerHTML = '';
    performance.clearResourceTimings()
    const promiseBuffer: Promise<void>[] = [];
    const resolution = "low";
    let totalDataTransfered = 0.0

    for (let i = 1; i <= TOTAL_IMAGES; i++) {
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
                totalDataTransfered += blob.size
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
    // Performance metrics
    const resources = performance.getEntriesByType("resource");
    let totalTime = 0;
    resources.forEach(entry => {
        if (entry.name.includes("/api/images")) {
            totalTime += entry.duration;
        }
    });

    const avgImgLoadTime = totalTime / resources.length;
    // Update performance metrics
    document.getElementById("totImgLoaded")!.textContent = `${resources.length}`
    document.getElementById("avgLoadTime")!.textContent = `${avgImgLoadTime.toFixed(3)} ms`;
    document.getElementById("data-transferred")!.textContent = `${(totalDataTransfered / 1024).toFixed(2)} KB`;
}


async function loadImage(imgElement: Element) {
    const displayImg = document.createElement("img");
    const resolution = "high";
    displayImg.width = 1280;
    displayImg.height = 720;
    performance.clearResourceTimings()
    let hiResImgDataTrans = 0.0;
    try {
        const res = await fetch(`/api/images/${resolution}/${imgElement.id}`);
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
        const blob = await res.blob()
        hiResImgDataTrans = blob.size;
        displayImg.src = URL.createObjectURL(blob)
    } catch (err) {
        console.error(err);
    }
    display.innerHTML = '';
    display.appendChild(displayImg);

    // Performance metrics

    const resources = performance.getEntriesByType("resource");
    let loadTime = 0;
    resources.forEach(entry => {
        loadTime = entry.duration;
    });

    // Update performance metrics
    document.getElementById("hiRes-imgLoadTime")!.textContent = `${loadTime.toFixed(3)} ms`;
    document.getElementById("data-transferred")!.textContent = `${(hiResImgDataTrans / (1024 * 1024)).toFixed(2)} MB`;
}


const display = document.getElementById("display-container")!;
const btn = document.getElementById("btn")!;
btn.addEventListener("click", loadPreviewImages);


export { loadImage, loadPreviewImages }
