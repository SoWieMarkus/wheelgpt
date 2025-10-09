declare global {
    namespace Express {
        export interface Request {
            channelId?: string;
            login?: string;
            legitTwitchSignature?: boolean;
        }
    }
}

export { }