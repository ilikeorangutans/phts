import { UploadQueueService } from './../../services/upload-queue.service';
import { Component, OnInit, Input } from '@angular/core';
import { PhotoService } from '../../services/photo.service';
import { Collection } from '../../models/collection';

@Component({
  selector: 'app-photo-upload',
  templateUrl: './photo-upload.component.html',
  styleUrls: ['./photo-upload.component.css']
})
export class PhotoUploadComponent implements OnInit {

  @Input() collection: Collection;

  dropFilesMessage = 'Drop files here...';

  message = this.dropFilesMessage;

  constructor(
    private photoService: PhotoService,
    private uploadQueue: UploadQueueService
  ) { }

  ngOnInit() {
  }

  getFiles(event: DragEvent): Array<File> {
    const result = Array<File>();
    const dt = event.dataTransfer;
    if (dt.items) {
      for (let i = 0; i < dt.items.length; i++) {
        const f = dt.items[i].getAsFile();

        result.push(f);
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
      this.uploadQueue.enqueue(this.collection,  file);
    }

    this.message = this.dropFilesMessage;

    return false;
  }

  onDragOver(event: DragEvent): Boolean {
    event.stopPropagation();
    event.preventDefault();

    const files = this.getFiles(event);
    this.message = `Upload ${files.length} files...`;

    return false;
  }

  onDragEnd(event: DragEvent): Boolean {
    return false;
  }
}
