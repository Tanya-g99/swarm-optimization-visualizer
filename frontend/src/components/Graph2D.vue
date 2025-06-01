<template>
    <div class="plot-container2D-wrapper card-shadow">
        <div ref="container2D" class="plot-container2D"></div>
        <div ref="heightScale" class="height-scale">
            <!-- Шкала высот -->
            <div class="height-scale">
                <div v-for="(_, index) in 11" :key="index" class="height-scale-label">
                    {{ (settings.minZ + (settings.maxZ - settings.minZ) * (index / 10)).toFixed(1) }}
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import * as THREE from "three";
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js';
import { CSS2DRenderer, CSS2DObject } from "three/examples/jsm/renderers/CSS2DRenderer.js";
import utils from 'utils/graph';

const props = defineProps({
    targetFunction: {
        type: Object,
        defauutilslt: null,
    },
    points: {
        default: [],
        type: Array,
    },
    xMin: {
        default: -5,
        type: Number,
    },
    xMax: {
        default: 5,
        type: Number,
    },
    yMin: {
        default: -5,
        type: Number,
    },
    yMax: {
        default: 5,
        type: Number,
    },
    settings: {
        type: Object,
        default: {
            minZ: 0,
            maxZ: 0,
            xScale: 1,
            yScale: 1,
            zScale: 1,
            segments: 300,
        }
    }
});

const container2D = ref(null);
const distance2D = ref(0);
let scene2D, camera2D, renderer2D, labelRenderer2D, controls2D;
let planeMesh2D, pointsGroup2D;
let lastPointCount = 0;

const pointSize = () => (1 + 20 / camera2D.zoom);

const initScene = () => {
    scene2D = new THREE.Scene();

    const aspect = container2D.value.clientWidth / container2D.value.clientHeight;
    const viewSize = 50;

    camera2D = new THREE.OrthographicCamera(
        -viewSize * aspect, viewSize * aspect,
        viewSize, -viewSize,
        0, 1000
    );

    // camera2D = new THREE.PerspectiveCamera(50, container2D.value.clientWidth / container2D.value.clientHeight, 0.1, 1000);

    renderer2D = new THREE.WebGLRenderer({ antialias: true, alpha: true });
    renderer2D.setSize(container2D.value.clientWidth, container2D.value.clientHeight);
    container2D.value.appendChild(renderer2D.domElement);

    labelRenderer2D = new CSS2DRenderer();
    labelRenderer2D.setSize(container2D.value.clientWidth, container2D.value.clientHeight);
    labelRenderer2D.domElement.style.position = "absolute";
    labelRenderer2D.domElement.style.top = "0px";

    labelRenderer2D.domElement.style.pointerEvents = "none";
    container2D.value.appendChild(labelRenderer2D.domElement);

    controls2D = new OrbitControls(camera2D, renderer2D.domElement);
    controls2D.enableDamping = true;
    controls2D.enableRotate = false;
    controls2D.enableZoom = true;
    controls2D.enablePan = true;

    distance2D.value = utils.getLabelStep(camera2D.position.z);
    drawFunction();
    distance2D.value = utils.getLabelStep(camera2D.position.z);
    drawAxes();
    drawPoints();
    animate();
};

const drawFunction = () => {
    if (props.targetFunction !== null) {
        if (planeMesh2D) scene2D.remove(planeMesh2D);

        const width = (props.xMax - props.xMin) / props.settings.xScale;
        const height = (props.yMax - props.yMin) / props.settings.yScale;

        const geometry = new THREE.PlaneGeometry(width, height, props.settings.segments, props.settings.segments);

        const planeCenterX = (props.xMax + props.xMin) / 2 / props.settings.xScale;
        const planeCenterY = (props.yMax + props.yMin) / 2 / props.settings.yScale;
        geometry.applyMatrix4(new THREE.Matrix4().makeTranslation(planeCenterX, planeCenterY, 0));

        const vertices = geometry.attributes.position.array;
        const colors = [];


        for (let i = 0; i < vertices.length; i += 3) {
            const z = props.targetFunction.evaluate({ x: vertices[i] * props.settings.xScale, y: vertices[i + 1] * props.settings.yScale });
            vertices[i + 2] = 0;
            if (isNaN(z)) {
                vertices[i + 1] = NaN;
                colors.push(0, 0, 0);
            } else {
                const color = new THREE.Color().setHSL(...utils.getColorForHeight(z, props.settings));
                colors.push(color.r, color.g, color.b);
            }

        }

        geometry.computeVertexNormals();

        geometry.setAttribute("color", new THREE.Float32BufferAttribute(colors, 3));

        const material = new THREE.MeshBasicMaterial({
            vertexColors: true,
            side: THREE.DoubleSide,
            alphaTest: 0.1, // убирает области без значений
        });

        planeMesh2D = new THREE.Mesh(geometry, material);

        const cameraZ = Math.max(width, height) * 1.3;

        camera2D.position.set(planeCenterX, planeCenterY, cameraZ);

        controls2D.target = new THREE.Vector3(planeCenterX, planeCenterY, 0);

        scene2D.add(planeMesh2D);
    }
};

const drawAxes = async () => {
    scene2D.children.slice().forEach(child => {
        if (child instanceof THREE.Line || child instanceof CSS2DObject) {
            scene2D.remove(child);
        }
    });

    let step = distance2D.value;
    const stepScale = 2 / 3;
    const h = 0.1;
    const offset = 0.18 * distance2D.value ** (2 / 3);

    const addAxisLabels = (axis, min, max, scale, label, color) => {
        const l = Math.min(Math.floor(min), 0),
            r = Math.max(Math.ceil(max), 0),
            st = step * scale;

        const materialX = new THREE.LineBasicMaterial({ color: color });
        const geometryX = new THREE.BufferGeometry().setFromPoints([
            new THREE.Vector3(...axis(l, scale, - stepScale * step, h)),
            new THREE.Vector3(...axis(r, scale, stepScale * step, h)),
        ]);
        const axisX = new THREE.Line(geometryX, materialX);
        scene2D.add(axisX);

        createLabel(label, new THREE.Vector3(...axis(r, scale, step, h)), label, 0);
        createLabel(label, new THREE.Vector3(...axis(l, scale, -step, h)), label, 0);


        for (let i = l; i <= r; i = Math.ceil(i + st)) {
            if (i != 0) createLabel(i.toString(), new THREE.Vector3(...axis(i, scale, 0, h)), label, offset);
        }
    };

    const getX = (val, scale, offset, h) => [val / scale + offset, 0, h];
    const getY = (val, scale, offset, h) => [0, val / scale + offset, h];

    addAxisLabels(getX, props.xMin, props.xMax, props.settings.xScale, "X", 0xff0000);
    addAxisLabels(getY, props.yMin, props.yMax, props.settings.yScale, "Y", 0x00ff00);

    createLabel('0'.toString(), new THREE.Vector3(0, offset, h), 'Y', offset);
};

const createLabel = (text, position, axis, offset) => {
    const label = new CSS2DObject(utils.createLabelDiv(text));

    label.position.set(
        position.x + (axis === 'Y' ? offset : 0),
        position.y + (axis === 'X' ? offset : 0),
        position.z);

    scene2D.add(label);
};

const drawPoints = () => { 
    if (!props.targetFunction) return;

    const numPoints = props.points.length;
    const dynamicScaleFactor = pointSize();

    if (pointsGroup2D && lastPointCount === numPoints) {
        for (let i = 0; i < numPoints; i++) {
            const [x, y] = props.points[i];
            const point = pointsGroup2D.children[i];
            point.position.set(x / props.settings.xScale, y / props.settings.yScale, 0);
            point.scale.set(dynamicScaleFactor, dynamicScaleFactor, dynamicScaleFactor);
        }
    } else {
        if (pointsGroup2D) {
            scene2D.remove(pointsGroup2D);
        }

        pointsGroup2D = new THREE.Group();
        scene2D.add(pointsGroup2D);

        props.points.forEach(([x, y]) => {
            const geometry = new THREE.SphereGeometry(0.02);
            const material = new THREE.MeshBasicMaterial({ color: 0xff0000 });
            const point = new THREE.Mesh(geometry, material);
            point.position.set(x / props.settings.xScale, y / props.settings.yScale, 0);
            point.scale.set(dynamicScaleFactor, dynamicScaleFactor, dynamicScaleFactor);

            pointsGroup2D.add(point);
        });

        lastPointCount = numPoints;
    }
};

const animate = () => {
    requestAnimationFrame(animate);
    controls2D.update();
    distance2D.value = utils.getLabelStep(100 / camera2D.zoom);

    const dynamicScaleFactor = pointSize();
    pointsGroup2D.children.forEach(point => {
        point.scale.set(dynamicScaleFactor, dynamicScaleFactor, dynamicScaleFactor);
    });

    renderer2D.render(scene2D, camera2D);
    labelRenderer2D.render(scene2D, camera2D);
};

onMounted(() => {
    initScene();
});

watch(() => [props.xMin, props.xMax, props.yMin, props.yMax], () => {
    drawFunction();
    drawAxes();
});

watch(() => [props.points], () => {
    drawPoints();
}, { deep: true });

watch(() => [props.targetFunction], () => {
    drawFunction();
}, { deep: true });

watch(() => props.settings.segments, () => {
    drawFunction();
})

watch(distance2D, () => {
    drawAxes();
})
</script>

<style lang="scss" scoped>
.plot-container2D-wrapper {
    display: flex;

    .plot-container2D {
        position: relative;
        height: 100%;
        margin: 0;
        flex: 1;
        border: 1px solid var(--color-card-border);
        background-color: var(--color-background-soft);
        overflow: hidden;
    }

    .height-scale {
        width: max-content;
        background: linear-gradient(to top,
                hsl(241.20000000000002, 100%, 50%) 0%,
                hsl(164.92586283673867, 100%, 50%) 10%,
                hsl(133.33208076541015, 100%, 50%) 20%,
                hsl(109.08931912975396, 100%, 50%) 30%,
                hsl(88.65172567347739, 100%, 50%) 40%,
                hsl(70.64584437780472, 100%, 50%) 50%,
                hsl(54.367283378954205, 100%, 50%) 60%,
                hsl(39.39760159998097, 100%, 50%) 70%,
                hsl(25.4641615308203, 100%, 50%) 80%,
                hsl(12.377588510216079, 100%, 50%) 90%,
                hsl(0, 100%, 50%) 100%);
        /* Плавный градиент от красного к синему */
        height: 100%;
        position: relative;
        display: flex;
        flex-direction: column-reverse;
        justify-content: space-between;
        align-items: center;
        color: var(--color-chart-line);
        border: 1px solid var(--color-card-border);
    }
}
</style>
