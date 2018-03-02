import { TestBed, inject } from '@angular/core/testing';

import { CollectionStoreService } from './collection-store.service';

describe('CollectionStoreService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [CollectionStoreService]
    });
  });

  it('should be created', inject([CollectionStoreService], (service: CollectionStoreService) => {
    expect(service).toBeTruthy();
  }));
});
