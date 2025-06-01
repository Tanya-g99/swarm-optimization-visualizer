<template>
    <div ref="container3D" class="plot-container3D card-shadow"> </div>
</template>

<script setup>
import { onMounted, watch, ref } from "vue";
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { CSS2DRenderer, CSS2DObject } from "three/examples/jsm/renderers/CSS2DRenderer.js";
import utils from 'utils/graph';

const props = defineProps({
    targetFunction: {
        type: Object,
        default: null,
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

const container3D = ref(null);
const distance3D = ref({
    center: 0,
    X: 0,
    Y: 0,
    Z: 0
});
let scene3D, camera3D, renderer3D, labelRenderer3D, controls3D;
let functionMesh3D, pointsGroup3D;

let lastPointCount = 0;
let axesGroup3D = new THREE.Group();

const updateDistance = () => {
    const c = { x: camera3D.position.x, y: camera3D.position.z, z: camera3D.position.y };
    distance3D.value.center = utils.getLabelStep(camera3D.position.length());
    distance3D.value.X = utils.getLabelStep(Math.sqrt(c.y ** 2 + c.z ** 2));
    distance3D.value.Y = utils.getLabelStep(Math.sqrt(c.x ** 2 + c.z ** 2));
    distance3D.value.Z = utils.getLabelStep(Math.sqrt(c.x ** 2 + c.y ** 2));
}

function calculatePointDistance(p) {
    return Math.sqrt((camera3D.position.x - p.position.x) ** 2 +
        (camera3D.position.z - p.position.y) ** 2 +
        (camera3D.position.y - p.position.z) ** 2);
}

const initScene = () => {
    scene3D = new THREE.Scene();
    const rotationMatrix = new THREE.Matrix4().makeRotationX(-Math.PI / 2);
    scene3D.applyMatrix4(rotationMatrix);
    camera3D = new THREE.PerspectiveCamera(50, container3D.value.clientWidth / container3D.value.clientHeight, 0.1, 1000);
    camera3D.position.set(5, 5, 5);

    renderer3D = new THREE.WebGLRenderer({ antialias: true, alpha: true });;
    renderer3D.setSize(container3D.value.clientWidth, container3D.value.clientHeight);
    container3D.value.appendChild(renderer3D.domElement);

    labelRenderer3D = new CSS2DRenderer();
    labelRenderer3D.setSize(container3D.value.clientWidth, container3D.value.clientHeight);
    labelRenderer3D.domElement.style.position = "absolute";
    labelRenderer3D.domElement.style.top = "0px";
    labelRenderer3D.domElement.style.pointerEvents = "none";
    container3D.value.appendChild(labelRenderer3D.domElement);

    controls3D = new OrbitControls(camera3D, renderer3D.domElement);
    controls3D.enableDamping = true;

    const light = new THREE.DirectionalLight(0xffffff, 1);
    light.position.set(1, 1, 110).normalize();
    scene3D.add(light);

    const light2 = new THREE.DirectionalLight(0xffffff, 1);
    light2.position.set(-10, -10, -10).normalize();
    scene3D.add(light2);

    const ambientLight = new THREE.AmbientLight(0xffffff, 1);
    scene3D.add(ambientLight);

    updateDistance();

    drawFunction();
    drawAxes();
    drawPoints();
    animate();
};

const drawFunction = () => {
    if (props.targetFunction !== null) {
        if (functionMesh3D) scene3D.remove(functionMesh3D);

        const width = (props.xMax - props.xMin) / props.settings.xScale;
        const height = (props.yMax - props.yMin) / props.settings.yScale;

        const geometry = new THREE.PlaneGeometry(width, height, props.settings.segments, props.settings.segments);

        // сдвигаем плоскость, чтобы ее центр совпал с центром системы
        const planeCenterX = (props.xMax + props.xMin) / props.settings.xScale / 2;
        const planeCenterY = (props.yMax + props.yMin) / props.settings.yScale / 2;
        geometry.applyMatrix4(new THREE.Matrix4().makeTranslation(planeCenterX, planeCenterY, 0));

        const vertices = geometry.attributes.position.array;
        const colors = [];

        for (let i = 0; i < vertices.length; i += 3) {
            const x = (vertices[i]) * props.settings.xScale;
            const y = (vertices[i + 1]) * props.settings.yScale;
            const z = props.targetFunction.evaluate({ x: x, y: y });
            if (isNaN(z)) {
                vertices[i + 2] = NaN;
                colors.push(0, 0, 0);
            } else {
                vertices[i + 2] = z / props.settings.zScale;
                const color = new THREE.Color().setHSL(...utils.getColorForHeight(z, props.settings));
                colors.push(color.r, color.g, color.b);
            }
        }

        geometry.computeVertexNormals(); // освещение

        geometry.setAttribute("color", new THREE.Float32BufferAttribute(colors, 3));

        const material = new THREE.MeshLambertMaterial({
            vertexColors: true,
            side: THREE.DoubleSide,
            transparent: true,
            opacity: 0.7,
        });

        functionMesh3D = new THREE.Mesh(geometry, material);
        scene3D.add(functionMesh3D);
    }
};

const updateFunction = () => {
    if (props.targetFunction !== null) {
        const width = (props.xMax - props.xMin) / props.settings.xScale;
        const height = (props.yMax - props.yMin) / props.settings.yScale;

        const geometry = new THREE.PlaneGeometry(width, height, props.settings.segments, props.settings.segments);
        const planeCenterX = (props.xMax + props.xMin) / props.settings.xScale / 2;
        const planeCenterY = (props.yMax + props.yMin) / props.settings.yScale / 2;
        geometry.applyMatrix4(new THREE.Matrix4().makeTranslation(planeCenterX, planeCenterY, 0));

        const vertices = geometry.attributes.position.array;
        const colors = [];

        for (let i = 0; i < vertices.length; i += 3) {
            const x = vertices[i] * props.settings.xScale;
            const y = vertices[i + 1] * props.settings.yScale;
            const z = props.targetFunction.evaluate({ x, y });
            if (isNaN(z)) {
                vertices[i + 2] = NaN;
                colors.push(0, 0, 0);
            } else {
                vertices[i + 2] = z / props.settings.zScale;
                const color = new THREE.Color().setHSL(...utils.getColorForHeight(z, props.settings));
                colors.push(color.r, color.g, color.b);
            }
        }

        geometry.computeVertexNormals();
        geometry.setAttribute("color", new THREE.Float32BufferAttribute(colors, 3));

        functionMesh3D.geometry = geometry;
    }
};

const drawPoints = () => {
    if (!props.targetFunction || !props.points) return;

    const numPoints = props.points.length;

    if (pointsGroup3D && lastPointCount === numPoints) {
        for (let i = 0; i < numPoints; i++) {
            const [x, y] = props.points[i];
            const z = props.targetFunction.evaluate({ x: x, y: y });
            const point = pointsGroup3D.children[i];

            point.position.set(
                x / props.settings.xScale,
                y / props.settings.yScale,
                z / props.settings.zScale
            );
        }
    } else {
        if (pointsGroup3D) {
            scene3D.remove(pointsGroup3D);
        }

        pointsGroup3D = new THREE.Group();
        scene3D.add(pointsGroup3D);

        props.points.forEach(([x, y]) => {
            const z = props.targetFunction.evaluate({ x: x, y: y });
            const geometry = new THREE.SphereGeometry(0.1);
            const material = new THREE.MeshBasicMaterial({ color: 0xff0000 });
            const point = new THREE.Mesh(geometry, material);
            point.position.set(
                x / props.settings.xScale,
                y / props.settings.yScale,
                z / props.settings.zScale);

            pointsGroup3D.add(point);
        });
        scene3D.add(pointsGroup3D);
        lastPointCount = numPoints;
    }
};

const drawAxes = () => {
    const steps = distance3D.value;
    if (axesGroup3D) {
        scene3D.remove(axesGroup3D);
    }
    axesGroup3D = new THREE.Group();

    const stepScale = 2 / 3;
    const getX = (val, scale, offset) => new THREE.Vector3(val / scale + offset, 0, 0);
    const getY = (val, scale, offset) => new THREE.Vector3(0, val / scale + offset, 0);
    const getZ = (val, scale, offset) => new THREE.Vector3(0, 0, val / scale + offset);

    const axisParams = [
        {
            color: 0xFF0000,  // X
            min: props.xMin,
            max: props.xMax,
            scale: props.settings.xScale,
            axisFunc: getX,
            direction: new THREE.Vector3(1, 0, 0),
        },
        {
            color: 0x00FF00,  // Y
            min: props.yMin,
            max: props.yMax,
            scale: props.settings.yScale,
            axisFunc: getY,
            direction: new THREE.Vector3(0, 1, 0),
        },
        {
            color: 0x0000FF,  // Z
            min: props.settings.minZ,
            max: props.settings.maxZ,
            scale: props.settings.zScale,
            axisFunc: getZ,
            direction: new THREE.Vector3(0, 0, 1),
        }
    ];

    axisParams.forEach(param => {
        const min = Math.min(Math.floor(param.min), 0),
            max = Math.max(Math.ceil(param.max), 0);
        const start = param.axisFunc(min, param.scale, -stepScale * steps.center);
        const end = param.axisFunc(max, param.scale, stepScale * steps.center);

        const material = new THREE.LineBasicMaterial({ color: param.color });
        const geometry = new THREE.BufferGeometry().setFromPoints([start, end]);
        const axis = new THREE.Line(geometry, material);
        axesGroup3D.add(axis);

        const arrowSize = 2 * stepScale * steps.center;
        const arrow = new THREE.ArrowHelper(param.axisFunc(max + 1, param.scale, 0), param.axisFunc(max, param.scale, -stepScale * steps.center), arrowSize, param.color);
        axesGroup3D.add(arrow);
    });

    scene3D.add(axesGroup3D);


    scene3D.children.slice().forEach(child => {
        if (child instanceof CSS2DObject) {
            scene3D.remove(child);
        }
    });

    const addAxisLabels = (axis, min, max, scale, step, label) => {
        createLabel(label, new THREE.Vector3(...axis(max, scale, steps.center)), label);
        createLabel(label, new THREE.Vector3(...axis(min, scale, -steps.center)), label);

        const l = Math.min(Math.floor(min), 0),
            r = Math.max(Math.ceil(max), 0),
            st = step * scale;
        for (let i = l; i <= r; i = Math.ceil(i + st)) {
            createLabel(i.toString(), new THREE.Vector3(...axis(i, scale, 0)), label);
        }
    };

    addAxisLabels(getX, props.xMin, props.xMax, props.settings.xScale, steps.X, 'X');
    addAxisLabels(getY, props.yMin, props.yMax, props.settings.yScale, steps.Y, 'Y');
    addAxisLabels(getZ, props.settings.minZ, props.settings.maxZ, props.settings.zScale, steps.Z, 'Z');
};

const createLabel = (text, position, axis) => {
    const label = new CSS2DObject(utils.createLabelDiv(text));
    label.position.copy(position);

    scene3D.add(label);
};

const animate = () => {
    requestAnimationFrame(animate);
    controls3D.update();

    updateDistance();

    pointsGroup3D.children.forEach(point => {
        const pointDistance = calculatePointDistance(point);
        const dynamicScaleFactor = 0.1 + pointDistance / 40;
        point.scale.set(dynamicScaleFactor, dynamicScaleFactor, dynamicScaleFactor);
    });

    renderer3D.render(scene3D, camera3D);
    labelRenderer3D.render(scene3D, camera3D);
};

onMounted(() => {
    initScene();
});

watch(() => [props.xMin, props.xMax, props.yMin, props.yMax], () => {
    if (pointsGroup3D) scene3D.remove(pointsGroup3D);
    drawAxes();
    updateFunction();
});

watch(() => props.points, () => {
    drawPoints();
}, { deep: true });

watch(() => props.targetFunction, () => {
    drawAxes();
    updateFunction();
});

watch(() => props.settings.segments, () => {
    updateFunction();
})

watch(distance3D, () => {
    drawAxes();

}, { deep: true })
</script>

<style lang="scss" scoped>
.plot-container3D {
    width: 100%;
    height: 100%;
    min-height: 500px;
    display: block;
    position: relative;
    border: 1px solid var(--color-card-border);
    background-color: var(--color-background-soft);

}
</style>