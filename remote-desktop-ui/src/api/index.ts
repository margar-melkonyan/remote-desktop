import Auth from "@/api/auth";
import User from "@/api/users";
import Session from "@/api/sessions";

class API {
  public api: object = {
    auth: new Auth(),
    users: new User(),
    sessions: new Session(),
  };
}

export default API;
