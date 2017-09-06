import { Injectable } from "@angular/core";
import { Collection } from "./collection";
import { HttpClient } from "@angular/common/http";
import 'rxjs/add/operator/toPromise';

@Injectable()
export class CollectionService {
    constructor(private http: HttpClient) { }

    collections: Collection[] = [];

    getRecent(count: number): Promise<Collection[]> {
        return Promise.resolve(this.collections);
    }

    save(collection: Collection): Promise<Collection> {
        console.log("posting");
        var result: Promise<Collection>;

        this.http
            .post("/api/collections", collection)
            .subscribe(
                data => {
                    result =  Promise.resolve(collection);
                },
                err => {
                    console.log("error");
                    console.log(err);
                    result = Promise.reject("");
                }
            );

        return result;
    }


}