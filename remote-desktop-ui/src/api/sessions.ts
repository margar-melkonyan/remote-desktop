import BaseAPI from "@/api/base";

class Session extends BaseAPI {
  protected URI: string = "api/v1/sessions";
  public urls: object = {
    index: (): string => `${this.baseURL}/${this.URI}/`,
    store: (): string => `${this.baseURL}/${this.URI}/`,
    edit: (id: number): string => `${this.baseURL}/${this.URI}/${id}/edit`,
    update: (id: number): string => `${this.baseURL}/${this.URI}/${id}`,
    delete: (id: number): string => `${this.baseURL}/${this.URI}/${id}`,
  };
}
export default Session;
