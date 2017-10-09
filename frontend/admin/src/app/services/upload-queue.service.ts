import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../models/collection';
import { PhotoService } from './photo.service';
import { Injectable } from '@angular/core';
import 'rxjs/add/operator/scan';
// import 'rxjs/add/operator/switchMap';

@Injectable()
export class UploadQueueService {

  private subject = new Subject<QueuedItem>();

  constructor(
    private photoService: PhotoService
  ) {
    this.subject.subscribe(item => {
      // TODO this is not serialized :|
      this.photoService.upload(item.collection, item.file);
    });
  }

  enqueue(collection: Collection, file: File) {
    this.subject.next(new QueuedItem(collection, file));
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
