import BaseAPI from "@/api/base";

class Room extends BaseAPI {
  protected URI: string = 'api/v1/rooms';
  public urls: object= {
    rooms: (): string => `${this.baseURL}/${this.URI}`,
    room: (id: number): string => `${this.baseURL}/${this.URI}/${id}`,
    my: (id: number): string => `${this.baseURL}/${this.URI}/my`,
    roomInfo: (id: number): string => `${this.baseURL}/${this.URI}/${id}/info`
  };
}
export default Room;
