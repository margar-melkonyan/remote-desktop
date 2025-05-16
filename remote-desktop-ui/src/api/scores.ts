import BaseAPI from "@/api/base";

class Score extends BaseAPI {
  protected URI: string = 'api/v1/scores';
  public urls: object= {
    scores: (): string => `${this.baseURL}/${this.URI}`,
  };
}
export default Score;
