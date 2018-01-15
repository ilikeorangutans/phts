import { Component, OnInit, Input, EventEmitter, Output } from '@angular/core';

import { Photo } from './../../models/photo';

@Component({
  selector: 'app-selectable-photo-container',
  templateUrl: './selectable-photo-container.component.html',
  styleUrls: ['./selectable-photo-container.component.css']
})
export class SelectablePhotoContainerComponent implements OnInit {

  @Input() photo: Photo;
  @Input() selected = false;

  @Output() change: EventEmitter<SelectedPhoto> =  new EventEmitter<SelectedPhoto>();

  constructor() { }

  ngOnInit() {
  }

  onChange(event) {
    const checked: boolean = event.target.checked;
    const selectedPhoto = new SelectedPhoto(checked, this.photo);
    this.change.emit(selectedPhoto);
  }
}

export class SelectedPhoto {
  constructor(
    readonly selected: boolean,
    readonly photo: Photo
  ) { }
}
