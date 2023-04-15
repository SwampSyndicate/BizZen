import { Service } from "./service";
import { User } from "./user";

export class Appointment {
    constructor(
        public service:Service,
        public user: User
    ) {}
}