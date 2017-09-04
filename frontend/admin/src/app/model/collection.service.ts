import { Injectable } from "@angular/core";
import { Collection } from "./collection";
import { HttpClient } from "@angular/common/http";

@Injectable()
export class CollectionService {
    constructor(private http: HttpClient) { }

    collections: Collection[] = [];

    getRecent(count: number): Promise<Collection[]> {
        return Promise.resolve(this.collections);
    }

    save(collection: Collection): void {
        console.log("posting");
        this.http
            .post("/api/collections", collection)
            .subscribe(
                data => {},
                err => {
                    console.log("error");
                    console.log(err);                    
                }
            );
    }

    
}