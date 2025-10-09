import axios from "axios";
import * as z from "zod";

// See: https://api2.mania.exchange/Method/Index/37
export const getTrackmaniaExchangeData = async (mapUid: string) => {
	const url = `https://trackmania.exchange/api/maps/get_map_info/uid/${mapUid}`;
	const result = await axios.get(url);
	if (result.data === "") return undefined;
	const TrackmaniaExchangeMapSchema = z.object({
		TrackID: z.number(),
	});

	const { success, data } = TrackmaniaExchangeMapSchema.safeParse(result.data);
	return success ? data : null;
};
