const container = document.getElementById("container")!;
const display = document.getElementById("display-container")!;

async function loadRandomImages() {
    for (let i = 1; i <= 11; i++) {
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
        container.appendChild(img);
    }
    const imgTags = document.querySelectorAll("#container img");
    imgTags.forEach(img => img.addEventListener("click", () => loadImage(img)));
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
btn.addEventListener("click", loadRandomImages);
