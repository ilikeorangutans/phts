import { Component, OnInit, Input } from '@angular/core';
import { PhotoService } from "../../services/photo.service";
import { Collection } from "../../models/collection";

@Component({
  selector: 'app-photo-upload',
  templateUrl: './photo-upload.component.html',
  styleUrls: ['./photo-upload.component.css']
})
export class PhotoUploadComponent implements OnInit {

  @Input() collection: Collection;

  constructor(
    private photoService: PhotoService
  ) { }

  ngOnInit() {
  }

  onDrop(event: DragEvent): Boolean {
    event.stopPropagation();
    event.preventDefault();
    console.log("onDrop()");

    let dt = event.dataTransfer;
    if (dt.items) {
      console.log("using data transfer item list")
      for (let i = 0; i < dt.items.length; i++) {
        let f = dt.items[i].getAsFile();
        console.log("got file", f);

        this.photoService.upload(this.collection, f);
      }
    } else {
      console.log("using data transfer")
      // TODO
    }



    return false;
  }

  onDragOver(event: DragEvent): Boolean {
    event.stopPropagation();
    event.preventDefault();

    return false;
  }

  onDragEnd(event: DragEvent): Boolean {
    console.log("onDragEnd()");
    return false;
  }
}
