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
    private photoService: PhotoService
  ) { }

  ngOnInit() {
  }

  getFiles(event: DragEvent): Array<File> {
    const result = Array<File>();
    const dt = event.dataTransfer;
    if (dt.items) {
      console.log('using data transfer item list');
      for (let i = 0; i < dt.items.length; i++) {
        const f = dt.items[i].getAsFile();
        console.log('got file', f);

        result.push(f);
      }
    } else {
      console.log('using data transfer');
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
    console.log('onDrop()');

    const dt = event.dataTransfer;
    if (dt.items) {
      console.log('using data transfer item list');
      for (let i = 0; i < dt.items.length; i++) {
        const f = dt.items[i].getAsFile();
        console.log('got file', f);

        // TODO would be cool if we had a queue
        this.photoService.upload(this.collection, f);
      }
    } else {
      console.log('using data transfer');
      // TODO
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
    console.log('onDragEnd()');
    return false;
  }
}
