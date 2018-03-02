import { TestBed, inject } from '@angular/core/testing';

import { PhotoStoreService } from './photo-store.service';

describe('PhotoStoreService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PhotoStoreService]
    });
  });

  it('should be created', inject([PhotoStoreService], (service: PhotoStoreService) => {
    expect(service).toBeTruthy();
  }));
});
