import express from "express";
import { AuthenticationController } from "../controllers";
import { requiresWebAuthentication } from "../middlewares/web-authentication";

const router = express.Router();

router.get("/twitch", AuthenticationController.twitchLoginRequest);
router.post("/login", AuthenticationController.login);
router.delete("/remove", requiresWebAuthentication, AuthenticationController.remove);
router.post("/renew", requiresWebAuthentication, AuthenticationController.updatePluginToken);
router.get("/token", requiresWebAuthentication, AuthenticationController.getPluginToken);

export default router;