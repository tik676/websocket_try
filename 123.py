import pygame

pygame.init()

WIDTH, HEIGHT = 500, 700
BACKGROUND_COLOR = (135, 206, 235)  
TREE_GREEN = (34, 139, 34)          
TRUNK_BROWN = (101, 67, 33)        
STAR_YELLOW = (255, 255, 0)         
ORNAMENT_RED = (255, 0, 0)          
ORNAMENT_BLUE = (0, 0, 255)         

screen = pygame.display.set_mode((WIDTH, HEIGHT))
pygame.display.set_caption("Новогодняя елка")

def draw_tree():
    trunk_width, trunk_height = 50, 80
    trunk_x = WIDTH // 2 - trunk_width // 2
    trunk_y = HEIGHT - trunk_height - 30
    pygame.draw.rect(screen, TRUNK_BROWN, (trunk_x, trunk_y, trunk_width, trunk_height))
    
    tree_base_y = trunk_y
    layer_heights = [100, 90, 80]  
    layer_widths = [200, 160, 120]  
    
    for i, (height, width) in enumerate(zip(layer_heights, layer_widths)):
        layer_y = tree_base_y - sum(layer_heights[:i+1]) + i * 20
        
   
        points = [
            (WIDTH // 2, layer_y),                    
            (WIDTH // 2 - width // 2, layer_y + height),  
            (WIDTH // 2 + width // 2, layer_y + height)   
        ]
        
        pygame.draw.polygon(screen, TREE_GREEN, points)
    
    star_y = tree_base_y - sum(layer_heights) - 40
    draw_star(WIDTH // 2, star_y, 20)
    
    draw_ornaments()

def draw_star(x, y, size):

    points = []
    import math
    for i in range(10):
        angle = math.pi * i / 5
        radius = size if i % 2 == 0 else size // 2
        point_x = x + radius * math.cos(angle - math.pi / 2)
        point_y = y + radius * math.sin(angle - math.pi / 2)
        points.append((point_x, point_y))
    
    pygame.draw.polygon(screen, STAR_YELLOW, points)

def draw_ornaments():

    ornament_positions = [
        (WIDTH // 2 - 60, HEIGHT - 200),
        (WIDTH // 2 + 40, HEIGHT - 180),
        (WIDTH // 2 - 30, HEIGHT - 280),
        (WIDTH // 2 + 60, HEIGHT - 300),
        (WIDTH // 2 - 80, HEIGHT - 350),
        (WIDTH // 2 + 20, HEIGHT - 380),
    ]
    

    for i, (x, y) in enumerate(ornament_positions):
        color = ORNAMENT_RED if i % 2 == 0 else ORNAMENT_BLUE
        pygame.draw.circle(screen, color, (x, y), 8)

running = True
clock = pygame.time.Clock()

while running:
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
    
 
    screen.fill(BACKGROUND_COLOR)
    
   
    draw_tree()
    

    pygame.display.flip()
    clock.tick(60)

pygame.quit()

