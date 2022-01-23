import { TestBed } from '@angular/core/testing';

import { WebsockService } from './websock.service';

describe('WebsockService', () => {
  let service: WebsockService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WebsockService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
