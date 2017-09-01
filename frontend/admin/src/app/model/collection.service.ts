import { Injectable } from "@angular/core";
import { Collection } from "./collection";

@Injectable()
export class CollectionService {
    collections: Collection[] = [];

    getRecent(count: number): Promise<Collection[]> {
        return Promise.resolve(this.collections);
    }

    
}