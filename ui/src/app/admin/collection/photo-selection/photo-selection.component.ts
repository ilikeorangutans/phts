import { Component, OnInit, Output, EventEmitter } from '@angular/core';

import { Photo } from './../../models/photo';
import { SelectedPhoto } from '../selectable-photo-container/selectable-photo-container.component';

@Component({
  selector: 'app-photo-selection',
  template: '',
  styleUrls: ['./photo-selection.component.css']
})
export class PhotoSelectionComponent implements OnInit {

  @Output() selectedPhotos: EventEmitter<Array<Photo>> = new EventEmitter<Array<Photo>>();

  selected: Array<Photo> = [];

  constructor() { }

  ngOnInit() {
  }

  onPhotoSelect(selected: SelectedPhoto) {
    console.log("onPhotoSelect()", selected);
    this.toggleSelectPhoto(selected.photo);
  }

  toggleSelectPhoto(photo: Photo) {
    if (this.selected.includes(photo)) {
      this.selected = this.selected.filter(p => p.id !== photo.id);
    } else {
      this.selected.push(photo);
    }
    this.selectedPhotos.emit(this.selected);
  }

  deselect() {
    this.selected = [];
    this.selectedPhotos.emit(this.selected);
  }

}
