import { Photo } from './../models/photo';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../models/collection';
import { PhotoService } from './photo.service';
import { Injectable } from '@angular/core';

import 'rxjs/add/operator/concatMap';
import 'rxjs/add/observable/fromPromise';

@Injectable()
export class UploadQueueService {

  private serializedQueue = new Subject<QueuedItem>();

  queue: BehaviorSubject<Array<File>> = new BehaviorSubject([]);

  successfulUploads: Subject<Photo> = new Subject();

  constructor(
    private photoService: PhotoService
  ) {
    this.serializedQueue
      .asObservable()
      .concatMap((item) => {
        return Observable
          .fromPromise(this.photoService.upload(item.collection, item.file));
      })
      .subscribe(photo => {
        const updated = this.queue.value;
        updated.shift();
        this.queue.next(updated);
        this.successfulUploads.next(photo);
      });
  }

  enqueue(collection: Collection, file: File) {
    this.serializedQueue.next(new QueuedItem(collection, file));
    const updated = this.queue.value;
    updated.push(file);
    this.queue.next(updated);
  }

}

class QueuedItem {
  collection: Collection;
  file: File;

  constructor(c: Collection, f: File) {
    this.collection = c;
    this.file = f;
  }
}
