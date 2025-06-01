function getLabelStep(distance) {
    let step = 0.1;

    if (distance > 5) step = 0.3;
    if (distance > 10) step = 0.7;
    if (distance > 15) step = 1;
    if (distance > 20) step = 2;
    if (distance > 50) step = 5;
    if (distance > 100) step = 10;
    if (distance > 300) step = 30;
    if (distance > 500) step = 50;
    if (distance > 700) step = 70;
    if (distance > 1000) step = 100;

    return step;
}

function getColorForHeight(z, settings) {
    const t = ((z - settings.minZ) / (settings.maxZ - settings.minZ)) ** 0.5;
    return [0.67 * (1 - t), 1, 0.5];
};

function createLabelDiv(text) {
    const div = document.createElement("div");
    div.textContent = text;

    div.style.color = "var(--color-chart-line)";
    div.style.fontSize = "14px";
    div.style.opacity = "1";

    return div;
}

export default {
    getLabelStep,
    getColorForHeight,
    createLabelDiv,
};