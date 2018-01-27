import { Collection } from './../../models/collection';
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { AlbumService } from '../../services/album.service';
import { Album } from '../../models/album';

@Component({
  selector: 'app-album-selector',
  templateUrl: './album-selector.component.html',
  styleUrls: ['./album-selector.component.css']
})
export class AlbumSelectorComponent implements OnInit {

  @Input() collection: Collection;

  @Output() albumSelected: EventEmitter<Album> = new EventEmitter<Album>();

  albums: Array<Album> = [];

  constructor(
    private albumService: AlbumService
  ) { }

  ngOnInit() {
    this.albumService.list(this.collection).then(albums => this.albums = albums);
  }

  onSelect(album: Album) {
    this.albumSelected.emit(album);
  }

}
