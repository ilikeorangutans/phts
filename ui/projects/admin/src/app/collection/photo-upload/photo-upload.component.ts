import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

import { Album } from './../../models/album';
import {
  UploadQueueService,
  UploadRequest,
} from './../../services/upload-queue.service';
import { Photo } from './../../models/photo';
import { Collection } from './../../models/collection';

@Component({
  selector: 'app-photo-upload',
  templateUrl: './photo-upload.component.html',
  styleUrls: ['./photo-upload.component.css'],
})
export class PhotoUploadComponent implements OnInit {
  @Input() collection: Collection;

  @Input() album: Album;

  dropFilesMessage = 'Drop files here...';

  message = this.dropFilesMessage;

  @Output() upload = new EventEmitter<Photo>();

  queue: Array<File> = new Array();

  constructor(private uploadQueue: UploadQueueService) {
    this.uploadQueue.successfulUploads.subscribe((photo) => {
      this.upload.emit(photo);
    });
  }

  ngOnInit() {}

  getFiles(event: DragEvent): Array<File> {
    const result = Array<File>();
    const dt = event.dataTransfer;
    if (dt === null) {
      return result;
    }

    const items = dt.items;
    if (items) {
      for (let i = 0; i < items.length; i++) {
        const f = items[i].getAsFile();

        if (f !== null) {
          result.push(f);
        }
      }
    } else {
      // TODO
    }

    return result;
  }

  onDragLeave(): Boolean {
    this.message = this.dropFilesMessage;
    return false;
  }

  onDrop(event: DragEvent): Boolean {
    event.stopPropagation();
    event.preventDefault();

    const files = this.getFiles(event);

    for (const file of files) {
      this.uploadQueue.enqueue(
        new UploadRequest(file, this.collection, this.album)
      );
    }

    this.message = this.dropFilesMessage;

    return false;
  }

  onDragOver(event: DragEvent): Boolean {
    event.stopPropagation();
    event.preventDefault();

    const fileCount = event.dataTransfer?.items.length;
    this.message = `Upload ${fileCount} files...`;

    return false;
  }

  onDragEnd(_: DragEvent): Boolean {
    return false;
  }

  filesSelected(event: Event) {
    const x = event.target;
    if (x === null) {
      return;
    }
    const files: FileList = x['files'];

    const result = Array<File>();
    for (let i = 0; i < files.length; i++) {
      const f = files.item(i);
      if (f === null) {
        continue;
      }
      result.push(f);
    }

    for (const file of result) {
      this.uploadQueue.enqueue(
        new UploadRequest(file, this.collection, this.album)
      );
    }
    if (x['value'] !== null) {
      x['value'] = '';
    }
  }
}
