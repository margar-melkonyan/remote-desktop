import Room from "@/api/rooms";
import Auth from "@/api/auth";
import User from "@/api/users"
import Score from "@/api/scores";

class API {
  public api: object = {
    rooms: new Room(),
    auth: new Auth(),
    users: new User(),
    scores: new Score(),
  }
}

export default API;
