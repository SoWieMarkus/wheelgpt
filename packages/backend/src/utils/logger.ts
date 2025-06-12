import { createLogger, format, transports } from "winston";

const logger = createLogger({
	level: "info",
	format: format.combine(
		format.timestamp(),
		format.printf(({ level, message, timestamp }) => {
			return `[${timestamp}] ${level}: ${message}`;
		}),
	),
	transports: [new transports.Console(), new transports.File({ filename: "logs/app.log" })],
});

export default logger;
