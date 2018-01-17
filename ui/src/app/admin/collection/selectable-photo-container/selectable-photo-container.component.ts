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

  @Output() change: EventEmitter<Photo> = new EventEmitter<Photo>();

  constructor() { }

  ngOnInit() {
  }

  onChange(event) {
    this.change.emit(this.photo);
  }
}

export class SelectedPhoto {
  constructor(
    readonly selected: boolean,
    readonly photo: Photo
  ) { }
}
