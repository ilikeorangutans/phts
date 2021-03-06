import { Component, OnInit, Input, EventEmitter, Output } from '@angular/core';

import { Photo } from './../../models/photo';

@Component({
  selector: 'photo-selectable-container',
  templateUrl: './selectable-photo-container.component.html',
  styleUrls: ['./selectable-photo-container.component.css'],
})
export class SelectablePhotoContainerComponent implements OnInit {
  @Input() photo: Photo;
  @Input() selected = false;

  @Output() change: EventEmitter<Photo> = new EventEmitter<Photo>();

  constructor() { }

  containerClasses() {
    return {
      'not-selected': !this.selected,
      'selected': this.selected,
    };
  }

  ngOnInit() { }

  classes() {
    return {
      fa: true,
      'fa-circle': !this.selected,
      'fa-check-circle': this.selected,
    };
  }

  onChange(event) {
    event.preventDefault();
    if (event.shiftKey) {
      // TODO build shift selection logic
    }
    this.change.emit(this.photo);

  }
}

export class SelectedPhoto {
  constructor(readonly selected: boolean, readonly photo: Photo) { }
}
