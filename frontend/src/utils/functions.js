const rastrigin = {
    label: "Функция Растригина",
    formula: "10 * 2 + (x^2 - 10*cos(2*PI*x)) + (y^2 - 10*cos(2*PI*y))"
};
const rosenbrock = {
    label: "Функция Розенброка",
    formula: "(1 - x)^2 + 100*(y - x^2)^2"
};
const ackley = {
    label: "Функция Акла",
    formula: "-20*exp(-0.2*sqrt(0.5*(x^2 + y^2))) - exp(0.5*(cos(2*PI*x) + cos(2*PI*y))) + exp(1) + 20"
};
const levi = {
    label: "Функция Леви",
    formula: "sin(3*PI*x)^2 + (x-1)^2 * (1+sin(3*PI*y)^2) + (y-1)^2 * (1+sin(2*PI*y)^2)"
};
const schwefel = {
    label: "Функция Швефеля",
    formula: "418.9829 * 2 - (x*sin(sqrt(abs(x))) + y*sin(sqrt(abs(y))))"
};

export default { rastrigin, rosenbrock, ackley, levi, schwefel };