import Auth from "@/api/auth";
import User from "@/api/users"

class API {
  public api: object = {
    auth: new Auth(),
    users: new User(),
  }
}

export default API;
