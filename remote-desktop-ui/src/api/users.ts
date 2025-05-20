import BaseAPI from "@/api/base";

class User extends BaseAPI {
  protected URI: string = "api/v1/users";
  public urls: object = {
    current: (): string => `${this.baseURL}/${this.URI}/current`,
  };
}
export default User;
