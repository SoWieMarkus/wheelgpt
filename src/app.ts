import cors from "cors";
import express, {
    json,
    type NextFunction,
    type Request,
    type Response,
} from "express";
import createHttpError, { isHttpError } from "http-errors";
import { Log } from "./utils";
import { AuthenticationRouter, GameClientRouter } from "./routes";

const app = express();

app.use(cors());
app.use(json());

app.use("/authentication", AuthenticationRouter);
app.use("/game", GameClientRouter);


// Handling of unknown endpoints
app.use((request, response, next) => {
    next(createHttpError(404, "Endpoint not found."));
});

// Error handling
app.use(
    (
        error: unknown,
        request: Request,
        response: Response,
        next: NextFunction,
    ) => {
        console.log(error);
        const errorMessage = isHttpError(error)
            ? error.message
            : "An unknown error occured.";
        const errorStatus = isHttpError(error) ? error.status : 500;
        Log.error(`Status ${errorStatus}: ${errorMessage}`);
        response.status(errorStatus).json({ error: errorMessage });
    },
);

export default app;
