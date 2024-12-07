import express from "express";
import { AuthenticationController } from "../controllers";
import { requiresWebAuthentication } from "../middleware/web-authentication";

const router = express.Router();

router.get("/twitch", AuthenticationController.twitch);
router.post("/login", AuthenticationController.authenticate);
router.delete("/remove", requiresWebAuthentication, AuthenticationController.remove);
router.post("/renew", requiresWebAuthentication, AuthenticationController.updateToken);
router.get("/token", requiresWebAuthentication, AuthenticationController.getToken);

export default router;