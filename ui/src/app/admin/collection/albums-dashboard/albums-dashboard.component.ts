import { RenditionConfigurationService } from './../../services/rendition-configuration.service';
import { Observable } from 'rxjs/Observable';
import { Collection } from './../../models/collection';
import { Subscription } from 'rxjs/Subscription';
import { CollectionService } from './../../services/collection.service';
import { AlbumService } from './../../services/album.service';
import { Album } from './../../models/album';
import { Component, OnInit, OnDestroy } from '@angular/core';

@Component({
  selector: 'app-albums-dashboard',
  templateUrl: './albums-dashboard.component.html',
  styleUrls: ['./albums-dashboard.component.css']
})
export class AlbumsDashboardComponent implements OnInit, OnDestroy {
  album: Album = new Album();
  albums: Array<Album> = [];
  collection: Collection;

  private sub: Subscription;

  constructor(
    private collectionService: CollectionService,
    private albumService: AlbumService,
    private renditionConfigService: RenditionConfigurationService
  ) { }

  ngOnInit() {
    this.sub = this.collectionService.current
      .switchMap(collection => {
        return this.renditionConfigService.forCollection(collection).map(configs => {
          collection.renditionConfigurations = configs;
          return collection;
        });
      })
      .switchMap(collection => {
        this.collection = collection;
        return this.albumService.list(collection);
      }).subscribe(albums => {
        this.albums = albums;
      });
  }

  listAlbums() {
    this.albumService.list(this.collection).then(albums => this.albums = albums);
  }

  onSubmit() {
    this.albumService.save(this.collection, this.album).then(a => {
      this.listAlbums();
    });
    this.album = new Album();
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

}
