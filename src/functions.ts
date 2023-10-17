import chalk from "chalk";

type colorType = "text" | "variable" | "error" //Establecemos los distintos tipos que podra detectar nuestra función

//Asignamos un color a cada una de las opciones
const themeColors = {
    text: "#1ff061",
    variable: "#0011ff",
    error: "#f70505"
}

// Establecemos un export de nuestra función 'color' lo que nos permitira utilizarla en otros archivos
export const color = (color: colorType, message: any) => {
    return chalk.hex(themeColors[color])(message)
}