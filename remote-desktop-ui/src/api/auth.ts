import BaseAPI from "@/api/base";

class Auth extends BaseAPI {
  protected URI: string = 'auth';
  public urls: object = {
    signIn: (): string => `${this.baseURL}/${this.URI}/sign-in`,
    signUp: (): string => `${this.baseURL}/${this.URI}/sign-up`
  };
}
export default Auth;
