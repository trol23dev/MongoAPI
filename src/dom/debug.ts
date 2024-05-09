import chalk from "chalk";

export default class Debug {

    private type: string;

    constructor(type: string) {
        this.type = type || "system";
    }

    print(message: any) {
        const debugMap = {
            system: { color: chalk.blue, emoji: '🔧'},
            warning: { color: chalk.yellow, emoji: '⚠️'},
            success: { color: chalk.green, emoji: '✅'},
            error: { color: chalk.red, emoji: '❌'},
        }
    
        if (!debugMap[this.type]) {
            this.type = 'system';
        }
    
        const { color, emoji } = debugMap[this.type];
        console.log(`[${emoji} | ${this.type.toUpperCase()}] ==> ${color(message)}`);
    }
}