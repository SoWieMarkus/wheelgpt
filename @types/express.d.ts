declare global {
    namespace Express {
        export interface Request {
            channelId?: string;
        }
    }
}

export { }