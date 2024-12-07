import axios from "axios";
import { z } from "zod";

export const getTrackmaniaExchangeId = async (mapUid: string) => {
    const url = `https://trackmania.exchange/api/maps/get_map_info/uid/${mapUid}`;
    const result = await axios.get(url);
    if (result.data === "") return undefined;
    const TrackmaniaExchangeMapSchema = z.object({
        TrackID: z.string()
    });
    const { success, data } = TrackmaniaExchangeMapSchema.safeParse(result.data);
    return success ? data.TrackID : undefined;
};